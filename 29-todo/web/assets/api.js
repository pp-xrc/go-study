/* 教学向：统一请求封装。token 放 localStorage，接口返回 {code,message,data} */

const TOKEN_KEY = "todo_token";
/* 空字符串 = 当前页面同源。Nginx 把 /user /todo 反代到 Gin :8089，Tunnel 也能用 */
const API_BASE = "";

function getToken() {
  return localStorage.getItem(TOKEN_KEY) || "";
}

function setToken(token) {
  localStorage.setItem(TOKEN_KEY, token);
}

function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

function goLogin() {
  if (!location.pathname.endsWith("/") && !location.pathname.endsWith("/index.html")) {
    location.href = "/";
  }
}

async function api(path, options) {
  const opts = options || {};
  const headers = Object.assign(
    { "Content-Type": "application/json" },
    opts.headers || {}
  );
  const token = getToken();
  if (token) {
    headers.Authorization = "Bearer " + token;
  }

  const res = await fetch(API_BASE + path, {
    method: opts.method || "GET",
    headers: headers,
    body: opts.body ? JSON.stringify(opts.body) : undefined,
  });

  let json;
  try {
    json = await res.json();
  } catch (e) {
    throw new Error("服务器返回不是 JSON");
  }

  if (json.code === 10005) {
    clearToken();
    goLogin();
    throw new Error(json.message || "请先登录");
  }
  if (json.code !== 0) {
    throw new Error(json.message || "请求失败");
  }
  return json.data;
}
