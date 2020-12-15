const HOST = window.location.origin;
const API_VERSION = "v1";

const request = ({ url, method = "get", body }) =>
  fetch(`${HOST}/${API_VERSION}/${url}`, {
    method,
    body: JSON.stringify(body),
    headers: {
      Authorization: `Bearer ${liff.getAccessToken()}`,
      "Content-Type": "application/json"
    }
  }).then(response => response.json());

export default {
  getTodos: () =>
    liff.getProfile().then(({ userId }) => request({ url: `todo/${userId}` })),
  getTodo: id =>
    liff
      .getProfile()
      .then(({ userId }) => request({ url: `todo/${userId}/${id}` })),
  editTodo: (id, { description, title }) =>
    liff.getProfile().then(({ userId }) =>
      request({
        url: `todo/${userId}/${id}`,
        method: "put",
        body: { description, title }
      })
    ),
  deleteTodo: id =>
    liff
      .getProfile()
      .then(({ userId }) =>
        request({ url: `todo/${userId}/${id}`, method: "delete" })
      ),
  createTodo: ({ description, title }) =>
    liff.getProfile().then(({ userId }) =>
      request({
        url: `todo/${userId}`,
        method: "post",
        body: { description, title }
      })
    )
};
