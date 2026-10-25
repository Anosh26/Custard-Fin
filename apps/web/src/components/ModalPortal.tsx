"use client";

import { useEffect, useState, ReactNode } from "react";
import { createPortal } from "react-dom";

interface ModalPortalProps {
  children: ReactNode;
}

/**
 * ModalPortal safely mounts modal dialogues directly into document.body.
 * This completely isolates modals from any parent CSS stacking contexts
 * (such as navbar backdrop-filter / transforms) that would otherwise trap
 * position:fixed elements and cause top clipping.
 */
export default function ModalPortal({ children }: ModalPortalProps) {
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    return () => setMounted(false);
  }, []);

  if (!mounted || typeof document === "undefined") {
    return null;
  }

  return createPortal(children, document.body);
}
