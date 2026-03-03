import JSZip from "jszip";
import { join, relative } from "path";
import { readdirSync, readFileSync, statSync, writeFileSync } from "fs";

export async function zipDir(sourceDir: string, outPath: string): Promise<void> {
  const zip = new JSZip();

  function addFolder(dir: string) {
    for (const file of readdirSync(dir)) {
      const fullPath = join(dir, file);
      const relativePath = relative(sourceDir, fullPath);
      if (statSync(fullPath).isDirectory()) {
        if (file === "node_modules") continue;
        addFolder(fullPath);
      } else {
        zip.file(relativePath, readFileSync(fullPath));
      }
    }
  }

  addFolder(sourceDir);
  const buffer = await zip.generateAsync({ type: "nodebuffer", compression: "DEFLATE" });
  writeFileSync(outPath, buffer);
}
