const dev = process.env.NODE_ENV === "development";

export const HTTPURL = dev ? "http://localhost:8000/" : "/";
export const WSURL = dev ? "ws://localhost:8001/ws" : "/ws";