import { expect, jest } from "@jest/globals";
import { UiServer } from "../server.mjs";
import fs from "fs";
import axios from "axios";

const PORT = 8000;
const TEST_GALLERY_DIR = "./test/test-gallery";
const TEST_ALBUM_ID = "test-album";

describe("Server", () => {
  let server;

  beforeAll(async () => {
    server = new UiServer(PORT, TEST_GALLERY_DIR);
    await server.start();
  });

  afterAll(async () => {
    await server.stop();
  });

  it("should return 200 for landing page", async () => {
    const res = await axios.get(`http://localhost:${PORT}`);
    expect(res.status).toEqual(200);
  });

  it("should return 200 for albums", async () => {
    const formData = new FormData();
    formData.append("myFile", fs.readFileSync("./test/photos/test1.jpg"));
    const res = await axios.post(
      `http://localhost:${PORT}/albums/album1/images`,
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
        body: formData,
      }
    );
    expect(res.status).toEqual(200);
    await new Promise((resolve) => setTimeout(resolve, 1000));
  });
});
