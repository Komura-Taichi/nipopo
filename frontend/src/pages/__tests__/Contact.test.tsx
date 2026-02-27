import { describe, test, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import Contact from "../Contact";

describe("Contact", () => {
    test("ページ見出し「問い合わせ先」が表示される", () => {
        render(<Contact />);
        expect(
            screen.getByRole("heading", { name: "問い合わせ先", level: 1 })
        ).toBeInTheDocument();
    });

    test("問い合わせページの説明文が表示される", () => {
        render(<Contact />);
        expect(
            screen.getByTestId("contact-description")
        ).toBeInTheDocument();
    })

    test("メールアドレス欄が正しく表示される", () => {
        render(<Contact />);
        expect(
            screen.getByText("メールアドレス")
        ).toBeInTheDocument();

        expect(
            screen.getByText(/.+\s*\[at\]\s*.+\..+/)
        ).toBeInTheDocument();
    });

    test("「SNS」の文字列が表示される", () => {
        render(<Contact />);
        expect(
            screen.getByText("SNS")
        ).toBeInTheDocument();
    });

    test("SNSリンクが機能しているか", () => {
        render(<Contact />);

        const github = screen.getByRole("link", { name: "GitHub" });
        const x = screen.getByRole("link", { name: "X" });

        expect(github).toHaveAttribute("href");
        expect(x).toHaveAttribute("href");
    });
});