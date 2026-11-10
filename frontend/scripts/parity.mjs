// Parity check (constitution I): the served UI must be the prototype, byte for byte.
// Compares prototype/index.html's <style> and <script> blocks with public/legacy/*,
// and the shell's <head> metadata with the prototype's. Exit 1 on any drift.
import { readFileSync } from "node:fs";
import { createHash } from "node:crypto";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";

const here = dirname(fileURLToPath(import.meta.url));
const root = join(here, "..");
const proto = readFileSync(join(root, "..", "prototype", "index.html"), "utf8");
const css = readFileSync(join(root, "public", "legacy", "app.css"), "utf8");
const js = readFileSync(join(root, "public", "legacy", "app.js"), "utf8");
const shell = readFileSync(join(root, "index.html"), "utf8");

const sha = (s) => createHash("sha256").update(s).digest("hex").slice(0, 12);
const block = (src, tag) => {
  const m = src.match(new RegExp(`<${tag}\\b[^>]*>([\\s\\S]*?)</${tag}>`));
  if (!m) throw new Error(`no <${tag}> in prototype`);
  return m[1];
};

const failures = [];
const check = (name, a, b) => {
  const ok = a === b;
  console.log(`${ok ? "ok  " : "FAIL"} ${name}  ${sha(a)} ${ok ? "==" : "!="} ${sha(b)}`);
  if (!ok) failures.push(name);
};

check("stylesheet == prototype <style>", css, block(proto, "style"));
check("bundle == prototype <script>", js, block(proto, "script"));

const heads = (s) => (s.match(/<(?:meta|title|link rel="icon")[^>]*>[^<]*(?:<\/title>)?/g) || []).filter((t) => !t.includes("stylesheet"));
check("head metadata (meta/title/icon)", heads(proto).join("\n"), heads(shell).join("\n"));
check("body root", block(proto, "body").includes('<div id="root"></div>') && block(shell, "body").includes('<div id="root"></div>') ? "root" : "missing", "root");

if (failures.length) {
  console.error(`\nParity drift in: ${failures.join(", ")}`);
  process.exit(1);
}
console.log("\nUI parity with prototype: OK");
