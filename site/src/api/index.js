import axios from "axios";
import { ElMessage } from 'element-plus';

const api = axios.create({
  baseURL: "/api/v1",
  timeout: 10000,
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

api.interceptors.response.use(
  (response) => {
    const res = response.data;
    if (res.code === 0) {
      return res.data;
    } else {
      // 404 等情况不弹错误提示，由组件自行处理
      if (res.code !== 404) {
        ElMessage.error(res.message || "请求失败")
      }
      return Promise.reject(new Error(res.message));
    }
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      window.location.href = "/login";
      return Promise.reject(error);
    }

    // 不在这里自动显示错误消息，让各组件自己处理
    // 这样可以避免重复弹窗，并且组件可以提供更有针对性的错误信息
    return Promise.reject(error);
  },
);

export const pollApi = {
  createPoll: (data) => api.post('/polls', data),
  getPoll: (id) => api.get(`/polls/${id}`),
  getPollByTopic: (topicId) => api.get(`/topics/${topicId}/poll`),
  submitVote: (data) => api.post('/polls/vote', data),
}

export const topicApi = {
  deleteTopic: (id) => api.delete(`/topics/${id}`),
}

export const postApi = {
  deletePost: (id) => api.delete(`/posts/${id}`),
}

export default api;
