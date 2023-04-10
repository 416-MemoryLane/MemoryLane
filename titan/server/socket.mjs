import { WebSocket, WebSocketServer } from "ws";
import { getAlbums } from "./app.mjs";

export const wss = new WebSocketServer({ noServer: true });

export const sendMessage = (type, message) => {
  wss.clients.forEach((client) => {
    if (client.readyState === WebSocket.OPEN) {
      client.send(JSON.stringify({ type, message }));
    }
  });
};

const handleMessage = ({ ws, data }) => {
  // receive messages from client here
  const message = JSON.parse(data);
};

wss.on("connection", (ws, req) => {
  sendMessage({
    message: "Connected to Memory Lane WebSocket Server",
    type: "connection",
  });

  sendMessage("albums", getAlbums());
  ws.on("message", (data) => {
    handleMessage({ ws, data });
  });
});
