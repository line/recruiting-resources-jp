import api from "./api";

function init(props, dom) {
  // Call API to get /todo/:id
  render(props, dom);
}

function render(props, dom) {
  if (props.id === "new" || props.id === "create") {
    switch (props.id) {
      case "create":
        api
          .createTodo({
            title: document.querySelector("#title").value,
            description: document.querySelector("#description").value
          })
          .then(todo => {
            window.location.href = "/#";
          });
        break;
      case "new":
        function create() {
          window.location.href = `${
            window.location.origin
          }/${window.location.hash
            .split("/")
            .slice(0, -1)
            .join("/")}/create`;
        }
        dom.innerHTML = `
        <div>
          <div>
            <label for='title'>Title: </label>
            <input for='title' id='title'/>
          </div>
          <div>
            <label for='description'>Description: </label>
            <textarea for='description' id='description'></textarea>
          </div>
          <button onclick='(${create})()'>Create</button>
        </div>
      `;
        break;
    }
  } else {
    switch (props.action) {
      case "save":
        const title = document.querySelector("#title").value;
        const description = document.querySelector("#description").value;
        api.editTodo(props.id, { title, description }).then(todo => {
          window.location.href = "/#";
        });
        break;
      case "edit":
        function save() {
          window.location.href = `${
            window.location.origin
          }/${window.location.hash
            .split("/")
            .slice(0, -1)
            .join("/")}/save`;
        }
        api.getTodo(props.id).then(todo => {
          dom.innerHTML = `
        <div>
          <div>
            <label for='title'>Title: </label>
            <input for='title' id='title' value='${todo.title}'/>
          </div>
          <div>
            <label for='description'>Description: </label>
            <textarea for='description' id='description'>${
              todo.description
            }</textarea>
          </div>
          <button onclick='(${save})()'>Save</button>
        </div>
      `;
        });
        break;
      case "delete":
        api
          .deleteTodo(props.id)
          .then(res => (window.location.href = window.location.origin + "/#"));
        break;
      default:
        api.getTodo(props.id).then(todo => {
          dom.innerHTML = `
        <div>
          <h1>${todo.title}</h1>
          <p>${todo.description}</p>
        </div>
      `;
        });
        break;
    }
  }
}

export default {
  init,
  render
};
