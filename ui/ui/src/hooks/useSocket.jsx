import { useEffect, useState } from "react";

export const useSocket = (setAlbums = () => null) => {
  const [socket, setSocket] = useState(null);
  const [isConnecting, setIsConnecting] = useState(false);

  useEffect(() => {
    if (!socket && !isConnecting) {
      setIsConnecting(true);
      return;
    }

    if (!socket && isConnecting) {
      const conn = new WebSocket("ws://localhost:4321");

      const ping = () => {
        setTimeout(() => {
          if (conn) {
            conn.send(JSON.stringify({ type: "ping" }));
            ping();
          }
        }, 5000);
      };

      conn.addEventListener("open", () => {
        setIsConnecting(false);
        setSocket(conn);
        // need to ping every so often to keep connection alive
        ping();
      });

      conn.addEventListener("message", (e) => {
        const data = JSON.parse(e.data);
        switch (data.type) {
          case "connection": {
            console.log(data.message);
            break;
          }
          case "albums": {
            console.log(data.message)
            setAlbums(data.message);
            break;
          }
          default:
            console.log(data);
        }
      });

      conn.addEventListener("close", () => {
        console.log("Socket connection closed");
        setSocket(null);
      });
    }
  }, [socket, isConnecting, setAlbums]);

  const send = (msg) => {
    socket?.send(JSON.stringify({ type: "message", message: msg }));
  };

  return { send, socket };
};
