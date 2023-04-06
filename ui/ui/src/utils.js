export const getEndpoint = (endpoint) => {
  if (process.env.NODE_ENV === "development") {
    return `http://localhost:4321${endpoint}`;
  } else {
    return endpoint;
  }
};
