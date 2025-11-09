import axios from "axios";
import { app_cfg } from "./app.cfg";

const endpoint = app_cfg.endpoint;

// crea un'istanza
const api = axios.create({
  baseURL: endpoint, // o la tua API base
});

// interceptor di risposta
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("username");
      // setUsername(null);
      // redirect globale
      window.location.href = "/"; // oppure "/login"
    }
    return Promise.reject(error);
  }
);

export default api;
