import ky from "ky";

export const httpClient = ky.create({
  timeout: 60000,
  hooks: {
    beforeRequest: [
      (req) => {
        if (!["GET", "HEAD", "OPTIONS"].includes(req.method)) {
          const csrfToken = getCsrfToken();
          if (csrfToken) {
            req.headers.set("x-csrf-token", csrfToken);
          }
        }
      },
    ],
  },
});

function getCsrfToken() {
  try {
    return getCookie("csrf");
  } catch {}
}

function refreshCsrfToken() {
  fetch("/api/auth/csrf", { method: "GET" });
}

setTimeout(refreshCsrfToken, 1000);

function getCookie(name: string): string | undefined {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop()!.split(";").shift();
}
