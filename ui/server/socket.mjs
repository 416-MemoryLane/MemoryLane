import { WebSocket, WebSocketServer } from "ws";

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

