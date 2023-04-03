import express from "express";
import fs from "fs";
import cors from "cors";
import multer from "multer";
import { v4 as uuidv4 } from "uuid";

const app = express();
app.use(express.json());
app.use(cors());

const port = 4321;
const ALBUM_DIR = "./albums";
const STATIC_PATH = "/static";

if (!fs.existsSync(ALBUM_DIR)) {
  fs.mkdirSync(ALBUM_DIR);
}

app.use(STATIC_PATH, express.static(ALBUM_DIR));

// all updates to FS should also update CRDT
// create album -> create crdt.json in the folder
// delete album -> delete entire directory?
// add/remove photos from album -> should update crdt
//  -> add: just add to added
//  -> delete: remove from added and add to deleted

const getFileExtension = (fileName) => {
  return fileName.split(".").at(-1);
};

const getPhotoUuid = (fileName) => {
  return fileName.split(".").at(0);
};

app.get("/albums", (req, res) => {
  const albums = fs.readdirSync(ALBUM_DIR, { withFileTypes: true });
  const albumStruct = albums
    .map((album) => {
      if (!album.isDirectory()) return null;
      const images = fs.readdirSync(`${ALBUM_DIR}/${album.name}`, {
        withFileTypes: true,
      });
      const crdtPath = `${ALBUM_DIR}/${album.name}/crdt.json`;
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
    .filter((album) => album !== null);
  res.status(200).send(albumStruct);
});

app.post("/albums", (req, res) => {
  const { albumName, uuid } = req.body;
  const crdt = {
    album: uuid,
    album_name: albumName,
    added: [],
    deleted: [],
  };
  const path = `${ALBUM_DIR}/${uuid}`;
  if (!fs.existsSync(path)) {
    fs.mkdirSync(path);
  }
  fs.writeFileSync(`${path}/crdt.json`, JSON.stringify(crdt));
  res.status(200).send("Album created");
});

app.delete("/albums/:uuid", (req, res) => {
  const { uuid } = req.params;
  const path = `${ALBUM_DIR}/${uuid}`;
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
  const albumPath = `${ALBUM_DIR}/${albumName}`;
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
  const albumPath = `${ALBUM_DIR}/${albumName}`;
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
    res.end("File is deleted successfully!");
  });
});

app.listen(port, () => {
  console.log("Server running on port " + port);
});
