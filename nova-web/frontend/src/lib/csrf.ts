function csrfToken() {
  return document.querySelector<HTMLMetaElement>('meta[name="csrf-token"]')?.content || "";
}

function csrfHeaders() {
  return {
    "X-CSRF-Token": csrfToken()
  };
}

export { csrfHeaders, csrfToken };
