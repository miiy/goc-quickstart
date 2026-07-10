import { useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

type PostCoverUploaderProps = DOMStringMap;

function PostCoverUploader(props: PostCoverUploaderProps) {
  const name = props.name || "cover";
  const label = props.label || "封面图";
  const accept = props.accept || "image/png,image/jpeg,image/webp";
  const inputRef = useRef<HTMLInputElement>(null);
  const [objectURL, setObjectURL] = useState("");
  const [fileMeta, setFileMeta] = useState<{ name: string; size: number } | null>(null);
  const previewURL = objectURL || props.existingUrl || "";

  useEffect(() => {
    return () => {
      if (objectURL) {
        URL.revokeObjectURL(objectURL);
      }
    };
  }, [objectURL]);

  return (
    <div className="space-y-2">
      <Label htmlFor={props.inputId || name}>
        {label}
      </Label>
      <div className="flex flex-col gap-2 sm:flex-row sm:items-start">
        <Input
          ref={inputRef}
          type="file"
          id={props.inputId || name}
          name={name}
          accept={accept}
          onChange={(event) => {
            const file = event.currentTarget.files?.[0];
            if (!file) {
              setObjectURL("");
              setFileMeta(null);
              return;
            }
            setFileMeta({ name: file.name, size: file.size });
            setObjectURL((current) => {
              if (current) {
                URL.revokeObjectURL(current);
              }
              return URL.createObjectURL(file);
            });
          }}
        />
        {fileMeta ? (
          <Button
            type="button"
            variant="outline"
            onClick={() => {
              if (inputRef.current) {
                inputRef.current.value = "";
              }
              setObjectURL((current) => {
                if (current) {
                  URL.revokeObjectURL(current);
                }
                return "";
              });
              setFileMeta(null);
            }}
          >
            清除
          </Button>
        ) : null}
      </div>
      <Card className="overflow-hidden bg-muted/40">
        <CardContent className="p-0">
          {previewURL ? (
            <img className="block max-h-60 w-full object-cover" src={previewURL} alt="" />
          ) : (
            <div className="flex min-h-24 items-center justify-center text-sm text-muted-foreground">未选择图片</div>
          )}
        </CardContent>
        {fileMeta ? (
          <CardFooter className="flex gap-4 px-3 py-2 text-sm text-muted-foreground">
            <span className="min-w-0 flex-1 truncate">{fileMeta.name}</span>
            <span>{formatBytes(fileMeta.size)}</span>
          </CardFooter>
        ) : null}
      </Card>
    </div>
  );
}

function formatBytes(bytes: number) {
  if (bytes < 1024) {
    return `${bytes} B`;
  }
  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }
  return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
}

export { PostCoverUploader };
