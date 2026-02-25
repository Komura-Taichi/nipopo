import { describe, expect, test, vi, beforeEach } from "vitest";
import { render, screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import "@testing-library/jest-dom";
import Home from "../Home";

const navigateMock = vi.fn();

vi.mock("react-router", async () => {
    const actual = await vi.importActual<typeof import("react-router")>("react-router");
    return {
        ...actual,
        useNavigate: () => navigateMock,
    };
});

// テストごとにnavigateの呼び出し履歴をリセット
beforeEach(() => {
    navigateMock.mockClear();
});

describe("Home", () => {
    test("タグの追加ができる", async () => {
        const user = userEvent.setup();

        render(<Home />);

        const tagInput = screen.getByPlaceholderText("タグを入力...");
        const tagAddButton = screen.getByRole("button", { name: "タグを追加" });

        await user.type(tagInput, "タグ1");
        await user.click(tagAddButton);

        const todayTagList = screen.getByLabelText("今日のタグ一覧")
        expect(within(todayTagList).getByText("タグ1")).toBeInTheDocument();
    });

    test("タグの削除ができる", async () => {
        const user = userEvent.setup();

        render(<Home />);

        const tagInput = screen.getByPlaceholderText("タグを入力...");
        const tagAddButton = screen.getByRole("button", { name: "タグを追加" });

        await user.type(tagInput, "タグ1");
        await user.click(tagAddButton);

        await user.click(screen.getByRole("button", { name: "タグ 「タグ1」 を削除" }));

        const todayTagList = screen.getByLabelText("今日のタグ一覧")
        expect(within(todayTagList).queryByText("タグ1")).not.toBeInTheDocument();
    });

    test("頑張り度で星を選ぶと正しく反映される", async () => {
        const user = userEvent.setup();

        render(<Home />);

        const starButton = screen.getByRole("radio", { "name": "3 / 5" });
        await user.click(starButton);
        expect(starButton).toBeChecked();

        expect(screen.getByLabelText("頑張り度の数字表記")).toHaveTextContent("3 / 5")
    });

    // TODO: APIアクセス時に有効化
    test.skip("できごと, タグ, 頑張り度を入力して記録ボタンを押すと記録ボタンが非アクティブになってテキストが変わる", async () => {
        const user = userEvent.setup();

        render(<Home />);

        const tagInput = screen.getByPlaceholderText("タグを入力...");
        const tagAddButton = screen.getByRole("button", { name: "タグを追加" });

        await user.type(tagInput, "タグ1");
        await user.click(tagAddButton);

        const contentInput = screen.getByPlaceholderText("今日一番心に残っていること");

        await user.type(contentInput, "今日はひたすらに研究をした。その結果、効率的な測定方法を知ることができた。");

        const starButton = screen.getByRole("radio", { name: "3 / 5" });
        await user.click(starButton);

        const saveButton = screen.getByRole("button", { name: "記録" });
        await user.click(saveButton);

        expect(saveButton).toBeDisabled();
        expect(saveButton).toHaveTextContent("記録中...");
    });

    test("できごとが抜けている状態で記録ボタンを押すとダイアログが表示される", async () => {
        const user = userEvent.setup();

        render(<Home />);

        const tagInput = screen.getByPlaceholderText("タグを入力...");
        const tagAddButton = screen.getByRole("button", { name: "タグを追加" });

        await user.type(tagInput, "タグ1");
        await user.click(tagAddButton);

        const starButton = screen.getByRole("radio", { name: "3 / 5" });
        await user.click(starButton);

        const saveButton = screen.getByRole("button", { name: "記録" });
        await user.click(saveButton);

        const invalidDialog = await screen.findByRole("dialog");
        expect(within(invalidDialog).getByRole("heading", { name: "タグ または 内容を入力してください" }));
    });
});