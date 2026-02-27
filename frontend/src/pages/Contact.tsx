import { FaGithub, FaXTwitter } from "react-icons/fa6";

import { borderStyle } from "../styles";

function Contact() {
    const githubUrl = "https://github.com/Komura-Taichi/nipopo";
    const xUrl = "https://x.com/buridge_stone?s=21";

    return (
        <div className="px-4 py-6">
            <h1 className="text-3xl font-bold">問い合わせ先</h1>

            <div className={`mt-5 ${borderStyle}`}>
                <p
                    className="text-md leading-relaxed"
                    data-testid="contact-description"
                >
                    不具合など発見した場合は、以下のメールアドレス、またはGitHubのIssue、Xにてご連絡ください。
                </p>

                <dl className="mt-8 space-y-8">
                    <div>
                        <dt className="text-xl font-semibold">メールアドレス</dt>
                        <dd className="mt-1 text-sm">
                            <span>tai4bsness [at] gmail.com</span>
                        </dd>
                    </div>

                    <div>
                        <dt className="text-xl font-semibold">SNS</dt>
                        <dd className="mt-2 flex items-center gap-3">
                            <a
                                href={githubUrl}
                                aria-label="GitHub"
                                target="_blank"
                                rel="noopener noreferer"
                                className="inline-flex items-center rounded p-1 hover:opacity-80 focus:outline-none focus:ring"
                            >
                                <FaGithub aria-hidden="true" focusable="false" size={40} />
                            </a>

                            <a
                                href={xUrl}
                                aria-label="X"
                                target="_blank"
                                rel="noopener noreferer"
                                className="inline-flex items-center rounded p-1 hover:opacity-80 focus:outline-none focus:ring"
                            >
                                <FaXTwitter aria-hidden="true" focusable="false" size={40} />
                            </a>
                        </dd>
                    </div>
                </dl>
            </div>
        </div>
    );
}

export default Contact;