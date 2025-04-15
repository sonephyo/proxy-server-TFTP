"use client";
import ThroughputVisualization from "@/component/visualization";
import { usePathname } from "next/navigation";

export default function Home() {
  const pathname = usePathname();

  // Get the last segment of the path, e.g., 'size4' from '/size4'
  const segments = pathname.split("/").filter(Boolean); // remove empty strings
  const filename = segments[segments.length - 1] || "";

  return (
    <div>
      <ThroughputVisualization filename={filename} />
    </div>
  );
}
