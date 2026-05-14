import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebars: SidebarsConfig = {
  docsSidebar: [
    "intro",
    "installation",
    "authentication",
    {
      type: "category",
      label: "TypeScript SDK — API Reference",
      items: [
        "typescript/stk-push",
        "typescript/c2b",
        "typescript/b2c",
        "typescript/b2b",
        "typescript/reversal",
        "typescript/transaction-status",
        "typescript/account-balance",
        "typescript/dynamic-qr",
      ],
    },
    {
      type: "category",
      label: "Python SDK — API Reference",
      items: [
        "python/stk-push",
        "python/c2b",
        "python/b2c",
        "python/b2b",
        "python/reversal",
        "python/transaction-status",
        "python/account-balance",
        "python/dynamic-qr",
      ],
    },
    {
      type: "category",
      label: "Go SDK — API Reference",
      items: [
        "go/stk-push",
        "go/c2b",
        "go/b2c",
        "go/b2b",
        "go/reversal",
        "go/transaction-status",
        "go/account-balance",
        "go/dynamic-qr",
      ],
    },
    "webhooks",
    {
      type: "category",
      label: "Advanced",
      items: [
        "typescript/retry-resilience",
        "typescript/framework-integrations",
      ],
    },
    "errors",
    "security",
    "production",
    {
      type: "category",
      label: "SDK Packages",
      items: [
        "typescript/index",
        "python/index",
        "go/index",
      ],
    },
    "faq",
  ],
};

export default sidebars;
