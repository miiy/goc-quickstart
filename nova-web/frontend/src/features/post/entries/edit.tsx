import { mountReactIslands, onDomReady } from "@/lib/reactIslands";

import { PostEditorPage } from "..";

onDomReady(() => {
  mountReactIslands("[data-post-editor-page]", PostEditorPage);
});
