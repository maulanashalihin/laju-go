import { describe, it, expect } from "vitest";
import { render } from "@testing-library/svelte";
import Button from "./Button.svelte";

describe("Button", () => {
  it("renders with default props", () => {
    const { container } = render(Button);
    const button = container.querySelector("button");
    expect(button).toBeTruthy();
    expect(button?.getAttribute("type")).toBe("button");
    expect(button?.getAttribute("disabled")).toBeNull();
  });

  it("renders as submit type", () => {
    const { container } = render(Button, { props: { type: "submit" } });
    const button = container.querySelector("button");
    expect(button?.getAttribute("type")).toBe("submit");
  });

  it("renders disabled state", () => {
    const { container } = render(Button, { props: { disabled: true } });
    const button = container.querySelector("button");
    expect(button?.getAttribute("disabled")).toBe("");
  });

  it("renders loading state", () => {
    const { container } = render(Button, { props: { loading: true } });
    const button = container.querySelector("button");
    expect(button?.getAttribute("disabled")).toBe("");
    expect(button?.innerHTML).toContain("svg");
  });

  it("applies variant classes", () => {
    const { container } = render(Button, { props: { variant: "danger" } });
    const button = container.querySelector("button");
    expect(button?.className).toContain("bg-red-600");
  });

  it("applies size classes", () => {
    const { container } = render(Button, { props: { size: "sm" } });
    const button = container.querySelector("button");
    expect(button?.className).toContain("px-3");
  });

  it("merges custom class names", () => {
    const { container } = render(Button, { props: { class: "custom-class" } });
    const button = container.querySelector("button");
    expect(button?.className).toContain("custom-class");
  });
});
