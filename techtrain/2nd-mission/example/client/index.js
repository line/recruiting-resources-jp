import topPage from "./topPage";
import todoDetailPage from "./todoDetailPage";

let app;

const routes = {
  "/": topPage,
  todo: todoDetailPage
};

function parseUrl() {
  const url = location.hash.slice(1).toLowerCase();
  const r = url.split("/");

  const request = {
    id: r[2],
    route: r[1] || "/",
    action: r[3]
  };

  return request;
}

function init() {
  app = document.getElementById("app");
  // Entry Point of LIFF App
  liff.init(
    {
      liffId: "YOUR LIFF ID"
    },
    () => {
      if (!liff.isInClient() && !liff.isLoggedIn()) {
        liff.login({
          redirectUri: window.location.origin +
            window.location.pathname
        });
      } else {
        router();
      }
    },
    err => console.error(err.code, err.message)
  );
}

function router() {
  const url = parseUrl();

  routes[url.route].init(url, app);
}

window.addEventListener("load", init);
window.addEventListener("hashchange", router);
