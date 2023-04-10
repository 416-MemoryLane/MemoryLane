import express from "express";
import fs from "fs";
import cors from "cors";
import multer from "multer";
import { v4 as uuidv4 } from "uuid";
import { sendMessage, wss } from "./socket.mjs";
import watch from "node-watch";
import getPort from "get-port";
import open from "open";
import {
  getFileExtension,
  getAlbums,
  getPhotoUuid,
  STATIC_PATH,
  GALLERY_DIR,
} from "./utils.mjs";
import * as dotenv from "dotenv";
import { SourceTextModule } from "vm";
dotenv.config({ path: "./.env" });

export class UiServer {
  constructor(port = 4321, galleryDir = GALLERY_DIR) {
    this.port = port;
    this.galleryDir = galleryDir;

    if (!fs.existsSync(this.galleryDir)) {
      fs.mkdirSync(this.galleryDir);
    }

    this.app = express();
    this.app.use(express.json());
    this.app.use(cors());

    this.watcher = watch(this.galleryDir, { recursive: true }, (_, name) => {
      if (!name.includes("crdt.json")) {
        sendMessage("albums", getAlbums(galleryDir));
      }
    });

    wss.on("connection", (ws, req) => {
      sendMessage({
        message: "Connected to Memory Lane WebSocket Server",
        type: "connection",
      });

      sendMessage("albums", getAlbums(galleryDir));
      ws.on("message", (data) => {
        // console.log(data)
      });
    });
  }

  async start() {
    this.app.use(STATIC_PATH, express.static(this.galleryDir));

    this.app.use("/", express.static("dist"));

    this.app.get("/login", (req, res) => {
      res.status(200).send({
        username: process.env.ML_USERNAME,
        password: process.env.ML_PASSWORD,
      });
    });

    this.app.get("/albums", (req, res) => {
      const albums = getAlbums(this.galleryDir);
      res.status(200).send(albums);
    });

    this.app.post("/albums", (req, res) => {
      const { albumName, uuid } = req.body;
      const crdt = {
        album: uuid,
        album_name: albumName,
        added: [],
        deleted: [],
      };
      const path = `${this.galleryDir}/${uuid}`;
      if (!fs.existsSync(path)) {
        fs.mkdirSync(path);
      }
      fs.writeFileSync(`${path}/crdt.json`, JSON.stringify(crdt));
      res.status(200).send("Album created");
    });

    this.app.delete("/albums/:uuid", (req, res) => {
      const { uuid } = req.params;
      const path = `${this.galleryDir}/${uuid}`;
      fs.rmSync(path, { recursive: true, force: true });
      res.status(200).send("Album deleted");
    });

    this.app.get("/albums/:id", (req, res) => {
      const { id } = req.params;
      console.log(id);
      res.send(id);
    });

    this.app.post("/albums/:albumName/images", (req, res) => {
      const { albumName } = req.params;
      const albumPath = `${this.galleryDir}/${albumName}`;
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
        res.status(200).send("File is uploaded successfully!");
      });
    });

    this.app.delete("/albums/:albumName/images/:imageName", (req, res) => {
      const { albumName, imageName } = req.params;
      const albumPath = `${this.galleryDir}/${albumName}`;
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

    if (["development", "test"].includes(process.env.NODE_ENV)) {
      this.server = this.app.listen(this.port, "0.0.0.0", () => {
        console.log("Server running on port " + this.port);
      });

      this.server.on("upgrade", (req, socket, head) => {
        wss.handleUpgrade(req, socket, head, (socket) => {
          wss.emit("connection", socket, req);
        });
      });
    } else {
      this.port = await getPort();
      // Start your server this.app.ication using the port number
      this.server = this.app.listen(this.port, "0.0.0.0", () => {
        console.log("Server running on port " + this.port);

        const uiUrl = `http://localhost:${this.port}/`;
        console.log(`Web UI accessible at ${uiUrl}`);
        return open(uiUrl);
      });

      this.server.on("upgrade", (req, socket, head) => {
        wss.handleUpgrade(req, socket, head, (socket) => {
          wss.emit("connection", socket, req);
        });
      });
    }
  }

  async stop() {
    wss.close();
    this.watcher.close();
    this.server.close();
  }
}
