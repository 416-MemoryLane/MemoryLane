import express from "express";
import fs from "fs";
import cors from "cors";
import multer from "multer";
import { v4 as uuidv4 } from "uuid";
import { sendMessage, wss } from "./socket.mjs";
import watch from "node-watch";
import getPort from "get-port";
import open from "open";
import * as dotenv from "dotenv";
dotenv.config({ path: "./.env" });

const app = express();
app.use(express.json());
app.use(cors());

const GALLERY_DIR = "../../memory-lane-gallery";
const STATIC_PATH = "/static";

if (!fs.existsSync(GALLERY_DIR)) {
  fs.mkdirSync(GALLERY_DIR);
}

watch(GALLERY_DIR, { recursive: true, delay: 750 }, (_, name) => {
  if (!name.includes("crdt.json")) {
    sendMessage("albums", getAlbums());
  }
});

app.use(STATIC_PATH, express.static(GALLERY_DIR));

app.use("/", express.static("dist"));

const getFileExtension = (fileName) => {
  return fileName.split(".").at(-1);
};

const getPhotoUuid = (fileName) => {
  return fileName.split(".").at(0);
};

export const getAlbums = () => {
  return (
    fs
      .readdirSync(GALLERY_DIR, { withFileTypes: true })
      .map((album) => {
        if (!album.isDirectory()) return null;
        const images = fs.readdirSync(`${GALLERY_DIR}/${album.name}`, {
          withFileTypes: true,
        });
        const crdtPath = `${GALLERY_DIR}/${album.name}/crdt.json`;
        if (!fs.existsSync(crdtPath)) {
          return res.status(500).send("No crdt.json file found in folder");
        }

        const crdt = JSON.parse(fs.readFileSync(crdtPath, "utf-8"));

        return {
          albumId: album.name,
          title: crdt.album_name,
          images: images
            .filter((file) => {
              const fileName = file.name;
              if (!fileName.includes(".")) {
                return false;
              }
              const extension = getFileExtension(fileName);
              if (
                !extension ||
                !["jpg", "png", "gif", "jpeg", "tiff"].includes(extension)
              ) {
                return false;
              }
              return file.isFile();
            })
            .map((file) => `${STATIC_PATH}/${album.name}/${file.name}`),
        };
      })
      .filter((album) => album !== null) ?? []
  );
};

app.get("/login", (req, res) => {
  res.status(200).send({
    username: process.env.ML_USERNAME,
    password: process.env.ML_PASSWORD,
  });
});

app.get("/albums", (req, res) => {
  const albums = getAlbums();
  res.status(200).send(albums);
});

app.post("/albums", (req, res) => {
  const { albumName, uuid } = req.body;
  const crdt = {
    album: uuid,
    album_name: albumName,
    added: [],
    deleted: [],
  };
  const path = `${GALLERY_DIR}/${uuid}`;
  if (!fs.existsSync(path)) {
    fs.mkdirSync(path);
  }
  fs.writeFileSync(`${path}/crdt.json`, JSON.stringify(crdt));
  res.status(200).send("Album created");
});

app.delete("/albums/:uuid", (req, res) => {
  const { uuid } = req.params;
  const path = `${GALLERY_DIR}/${uuid}`;
  fs.rmSync(path, { recursive: true, force: true });
  res.status(200).send("Album deleted");
});

app.get("/albums/:id", (req, res) => {
  const { id } = req.params;
  console.log(id);
  res.send(id);
});

app.post("/albums/:albumName/images", (req, res) => {
  const { albumName } = req.params;
  const albumPath = `${GALLERY_DIR}/${albumName}`;
  if (!fs.existsSync(albumPath)) {
    fs.mkdirSync(albumPath);
  }

  const photoUuid = uuidv4();
  const storage = multer.diskStorage({
    destination: function (req, file, callback) {
      callback(null, albumPath);
    },
    filename: function (req, file, callback) {
      callback(null, `${photoUuid}.${getFileExtension(file.originalname)}`);
    },
  });
  const upload = multer({ storage: storage }).single("myfile");
  upload(req, res, function (err) {
    if (err) {
      return res.end("Error uploading file.");
    }

    const crdtPath = `${albumPath}/crdt.json`;
    if (!fs.existsSync(crdtPath)) {
      return res
        .status(500)
        .send("No crdt.json file found for album. Aborting");
    }
    const crdt = JSON.parse(fs.readFileSync(crdtPath, "utf-8"));
    crdt.added.push(photoUuid);
    fs.writeFileSync(crdtPath, JSON.stringify(crdt));
    res.end("File is uploaded successfully!");
  });
});

app.delete("/albums/:albumName/images/:imageName", (req, res) => {
  const { albumName, imageName } = req.params;
  const albumPath = `${GALLERY_DIR}/${albumName}`;
  const imagePath = `${albumPath}/${imageName}`;

  fs.unlink(imagePath, function (err) {
    if (err) {
      return res.end("Error deleting file.");
    }
    const crdtPath = `${albumPath}/crdt.json`;
    if (!fs.existsSync(crdtPath)) {
      return res
        .status(500)
        .send("No crdt.json file found for album. Aborting");
    }
    const crdt = JSON.parse(fs.readFileSync(crdtPath, "utf-8"));
    const photoUuid = getPhotoUuid(imageName);
    crdt.added = crdt.added.filter((uuid) => uuid !== photoUuid);
    crdt.deleted.push(photoUuid);
    fs.writeFileSync(crdtPath, JSON.stringify(crdt));
    sendMessage("albums", getAlbums());
    res.end("File is deleted successfully!");
  });
});

if (process.env.NODE_ENV === "development") {
  const server = app.listen(4321, "0.0.0.0", () => {
    console.log("Server running on port " + 4321);
  });

  server.on("upgrade", (req, socket, head) => {
    wss.handleUpgrade(req, socket, head, (socket) => {
      wss.emit("connection", socket, req);
    });
  });
} else {
  getPort().then((port) => {
    // Start your server application using the port number
    const server = app.listen(port, "0.0.0.0", () => {
      console.log("Server running on port " + port);

      const uiUrl = `http://localhost:${port}/`;
      console.log(`Web UI accessible at ${uiUrl}`);
      return open(uiUrl);
    });

    server.on("upgrade", (req, socket, head) => {
      wss.handleUpgrade(req, socket, head, (socket) => {
        wss.emit("connection", socket, req);
      });
    });
  });
}
