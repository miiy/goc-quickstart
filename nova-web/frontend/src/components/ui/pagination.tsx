import * as React from "react";

import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

function Pagination({ className, ...props }: React.ComponentProps<"nav">) {
  return (
    <nav
      aria-label="pagination"
      data-slot="pagination"
      className={cn("mx-auto flex w-full justify-center", className)}
      {...props}
    />
  );
}

function PaginationContent({ className, ...props }: React.ComponentProps<"ul">) {
  return <ul data-slot="pagination-content" className={cn("flex flex-row items-center gap-1", className)} {...props} />;
}

function PaginationItem({ ...props }: React.ComponentProps<"li">) {
  return <li data-slot="pagination-item" {...props} />;
}

function PaginationLink({
  className,
  isActive,
  disabled,
  ...props
}: React.ComponentProps<"a"> & {
  isActive?: boolean;
  disabled?: boolean;
}) {
  return (
    <a
      aria-current={isActive ? "page" : undefined}
      aria-disabled={disabled ? "true" : undefined}
      data-active={isActive}
      data-slot="pagination-link"
      tabIndex={disabled ? -1 : props.tabIndex}
      className={cn(
        buttonVariants({
          variant: isActive ? "outline" : "ghost",
          size: "icon"
        }),
        "size-9 cursor-pointer",
        disabled && "pointer-events-none opacity-50",
        className
      )}
      {...props}
    />
  );
}

function PaginationPrevious({ className, ...props }: React.ComponentProps<typeof PaginationLink>) {
  return (
    <PaginationLink aria-label="上一页" className={cn("w-auto gap-1 px-2.5", className)} {...props}>
      <span aria-hidden="true">‹</span>
      <span>上一页</span>
    </PaginationLink>
  );
}

function PaginationNext({ className, ...props }: React.ComponentProps<typeof PaginationLink>) {
  return (
    <PaginationLink aria-label="下一页" className={cn("w-auto gap-1 px-2.5", className)} {...props}>
      <span>下一页</span>
      <span aria-hidden="true">›</span>
    </PaginationLink>
  );
}

function PaginationEllipsis({ className, ...props }: React.ComponentProps<"span">) {
  return (
    <span
      aria-hidden="true"
      data-slot="pagination-ellipsis"
      className={cn("flex size-9 items-center justify-center", className)}
      {...props}
    >
      ...
    </span>
  );
}

export {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious
};
