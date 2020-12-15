import api from "./api";

function init(props, dom) {
  api.getTodos().then(data => {
    props.todos = data.list;
    render(props, dom);
  });

  render(props, dom);
}

function render(props, dom) {
  let dataToRender = `<div>Loading...</div>`;
  const todos = props.todos || [];

  if (todos.length) {
    dataToRender = `
    <ul>${todos.reduce((acc, todo) => {
      return (acc += `<li><a href='/#/todo/${todo.id}'>${
        todo.title
      }</a> - <a href='/#/todo/${todo.id}/edit'>Edit</a> | <a href='/#/todo/${
        todo.id
      }/delete'>Delete</a></li>`);
    }, "")}</ul>
    `;
  } else {
    dataToRender = `<div>You have no ToDo's</div>`;
  }

  dataToRender += `
  <div><a href='/#/todo/new'>Create a New Todo</a></div>`;

  dom.innerHTML = dataToRender;
}

export default {
  init,
  render
};
