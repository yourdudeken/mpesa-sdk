import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebars: SidebarsConfig = {
  docsSidebar: [
    "intro",
    "installation",
    {
      type: "category",
      label: "TypeScript SDK (@yourdudeken/mpesa-sdk)",
      link: { type: "doc", id: "typescript/index" },
      items: [
        "typescript/index",
      ],
    },
    {
      type: "category",
      label: "Python SDK (yourdudeken-mpesa-sdk)",
      link: { type: "doc", id: "python/index" },
      items: [
        "python/index",
      ],
    },
    {
      type: "category",
      label: "Go SDK (github.com/yourdudeken/mpesa-sdk)",
      link: { type: "doc", id: "go/index" },
      items: [
        "go/index",
      ],
    },
    "authentication",
    "webhooks",
    "errors",
    "security",
    "production",
    "faq",
  ],
};

export default sidebars;
