import "@testing-library/jest-dom";

if (typeof HTMLDialogElement !== "undefined") {
    HTMLDialogElement.prototype.showModal ??= function() {
        // showModalが呼ばれたら開いたことにする
        this.setAttribute("open", "");
    };

    HTMLDialogElement.prototype.close ??= function() {
        // closeが呼ばれたら閉じたことにする
        this.removeAttribute("open");
    };
}