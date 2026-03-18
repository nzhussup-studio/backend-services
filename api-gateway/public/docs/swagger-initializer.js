window.onload = function () {
  window.ui = SwaggerUIBundle({
    dom_id: "#swagger-ui",
    urls: [
      { url: "/v1/base/v3/api-docs", name: "Base Service" },
      { url: "/v1/user/v3/api-docs", name: "User Service" },
      { url: "/auth/v3/api-docs", name: "Auth Service" },
      { url: "/v1/album/docs/doc.json", name: "Image Service" },
      { url: "/v1/llm/docs/doc.json", name: "LLM Service" },
    ],
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    layout: "StandaloneLayout",
  });
};
