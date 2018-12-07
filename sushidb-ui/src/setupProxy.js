const proxy = require("http-proxy-middleware");

// https://facebook.github.io/create-react-app/docs/proxying-api-requests-in-development
module.exports = function(app) {
  app.use(
    proxy("/api", {
      target: "http://localhost:3000/",
      pathRewrite: { "^/api": "" }
    })
  );
};
