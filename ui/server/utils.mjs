import fs from "fs";

export const STATIC_PATH = "/static";
export const GALLERY_DIR = "../../memory-lane-gallery";

export const getAlbums = (galleryDir = GALLERY_DIR) => {
  return fs
    .readdirSync(galleryDir, { withFileTypes: true })
    .map((album) => {
      if (!album.isDirectory()) return [];

      const images = fs.readdirSync(`${galleryDir}/${album.name}`, {
        withFileTypes: true,
      });
      const crdtPath = `${galleryDir}/${album.name}/crdt.json`;

      if (!fs.existsSync(crdtPath)) {
        throw new Error("No crdt.json file found in folder");
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
};

export const getFileExtension = (fileName) => fileName.split(".").at(-1);

export const getPhotoUuid = (fileName) => fileName.split(".").at(0);
