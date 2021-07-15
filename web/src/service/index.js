import {createApi} from "../util";

export const genApi = data => createApi("/api/struct/generate", data);
