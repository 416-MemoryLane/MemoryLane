import { useState, useEffect } from "react";
import { MinusIcon } from "./Icons/MinusIcon";
import { CloseIcon } from "./Icons/CloseIcon";
import { PhotoAlbum } from "react-photo-album";
import { Modal } from "./Modal";
import Lightbox from "yet-another-react-lightbox";

import "yet-another-react-lightbox/styles.css";
import Fullscreen from "yet-another-react-lightbox/plugins/fullscreen";
import Slideshow from "yet-another-react-lightbox/plugins/slideshow";
import Thumbnails from "yet-another-react-lightbox/plugins/thumbnails";
import "yet-another-react-lightbox/plugins/thumbnails.css";
import { getEndpoint } from "./utils";

export const AlbumOverview = ({ albumId, albumTitle, images, onBackClick }) => {
  const [index, setIndex] = useState(-1);
  const [thumbnails, setThumbnails] = useState([]);
  const [photos, setPhotos] = useState([]);
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);

  const handleFileUpload = async () => {
    const input = document.createElement("input");
    input.type = "file";
    input.accept = "image/*";
    input.onchange = async (event) => {
      const file = event.target.files[0];
      const formData = new FormData();
      formData.append("myfile", file);
      try {
        const response = await fetch(getEndpoint(`/albums/${albumId}/images`), {
          method: "POST",
          body: formData,
        });
        const data = await response.text();
        console.log(data);
      } catch (error) {
        console.error(error);
      }
    };
    input.click();
  };

  useEffect(() => {
    const promises = images.map((image) => {
      return new Promise((resolve, reject) => {
        const img = new Image();
        img.onload = () => {
          const width = img.width;
          const height = img.height;
          resolve({
            src: image,
            width: width,
            height: height,
            fileName: image.split("/").pop(),
          });
        };
        img.onerror = () => {
          reject(new Error(`Failed to load image: ${image}`));
        };
        img.src = getEndpoint(image);
      });
    });

    Promise.all(promises)
      .then((results) => {
        const thumbnails = results.map((result) => {
          return {
            src: getEndpoint(result.src),
            width: result.width,
            height: result.height,
            fileName: result.fileName,
          };
        });
        setThumbnails(thumbnails);
        setPhotos(
          results.map((result) => {
            return { src: getEndpoint(result.src) };
          })
        );
      })
      .catch((error) => {
        console.error(error);
      });
  }, [images]);

  const handleFileDelete = async (fileName) => {
    try {
      const response = await fetch(
        getEndpoint(`/albums/${albumId}/images/${fileName}`),
        {
          method: "DELETE",
        }
      );
      const data = await response.text();
      console.log(data);
    } catch (error) {
      console.error(error);
    }
  };

  const handleDeleteAlbum = async () => {
    try {
      const user = localStorage.getItem("galactus-user");
      const token = localStorage.getItem("galactus-token");
      const galactusResponse = await fetch(
        `https://memory-lane-381119.wl.r.appspot.com/delete_album?uuid=${albumId}&username=${user}`,
        {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      if (galactusResponse.ok) {
        const uiServerResponse = await fetch(
          getEndpoint(`/albums/${albumId}`),
          {
            method: "DELETE",
          }
        );
        if (uiServerResponse.ok) {
          onBackClick();
        }
      }
      const data = await galactusResponse.text();
      console.log(data);
    } catch (error) {
      console.error(error);
    }
  };

  const albumContainer = ({ containerProps, children, containerRef }) => (
    <div
      style={{
        border: "2px solid #ccc",
        borderRadius: "10px",
        padding: "20px",
      }}
    >
      <div ref={containerRef} {...containerProps}>
        {children}
      </div>
    </div>
  );

  const photoFrame = ({
    layout,
    layoutOptions,
    imageProps: { style, ...restImageProps },
    photo,
  }) => (
    <div
      style={{
        border: "2px solid #eee",
        borderRadius: "4px",
        boxSizing: "content-box",
        alignItems: "center",
        width: style?.width,
        maxWidth: "50%",
        padding: `20px`,
        margin: "5px",
        position: "relative",
      }}
    >
      <img style={{ ...style, width: "100%" }} {...restImageProps} />
      <MinusIcon
        className="text-gray-600 hover:fill-current hover:text-sky-400 cursor-pointer absolute top-0 right-0 opacity-50 px-1"
        onClick={() => handleFileDelete(photo.fileName)}
      />
    </div>
  );

  return (
    <div>
      <Modal
        isOpen={deleteModalOpen}
        callbackFn={() => setDeleteModalOpen(false)}
      >
        <div className="flex flex-col">
          <p className="text-2xl">
            Are you sure you want to delete this album?
          </p>
          <div className="flex flex-row self-center gap-3 mt-7">
            <p
              className="px-8 py-2 bg-green-200 border-green-300 border-[2px] rounded-md hover:cursor-pointer"
              onClick={handleDeleteAlbum}
            >
              Yes
            </p>
            <p
              className="px-8 py-2 bg-gray-200 border-gray-300 border-[2px] rounded-md hover:cursor-pointer"
              onClick={() => setDeleteModalOpen(false)}
            >
              No
            </p>
          </div>
        </div>
      </Modal>
      <div className="mx-auto px-8 sm:px-10 lg:px-12 lg:py-3">
        <div className="flex justify-between items-center">
          <div>
            <button
              className="px-3 hover:text-sky-400"
              onClick={handleFileUpload}
            >
              <span className="text-md font-semibold">Add Photo</span>
            </button>
            <button
              className="px-3 hover:text-red-500"
              onClick={() => setDeleteModalOpen(true)}
            >
              <span className="text-md font-semibold">Delete Album</span>
            </button>
          </div>
          <h1 className="text-2xl font-bold">{albumTitle}</h1>
          <CloseIcon
            className="hover:fill-current hover:text-sky-400 cursor-pointer"
            onClick={onBackClick}
          ></CloseIcon>
        </div>
        <PhotoAlbum
          photos={thumbnails}
          layout="rows"
          spacing={5}
          onClick={({ index }) => setIndex(index)}
          renderContainer={albumContainer}
          renderPhoto={(props) => photoFrame({ ...props })}
        />
        <Lightbox
          slides={photos}
          open={index >= 0}
          index={index}
          close={() => setIndex(-1)}
          on={{
            view: ({ index }) => {
              setIndex(index);
            },
          }}
          plugins={[Fullscreen, Slideshow, Thumbnails]}
        />
      </div>
    </div>
  );
};
