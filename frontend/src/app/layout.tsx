import type { Metadata } from "next";
import { JetBrains_Mono } from "next/font/google";
import "./globals.css";
import { Toaster } from "@/components/ui/sonner";

const jetbrains = JetBrains_Mono({ weight: "variable", subsets: ["cyrillic"] });

export const metadata: Metadata = {
  title: "Noob Chat",
  description: "Simple and Cool Chat Application",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={jetbrains.className}>{children}</body>
      <Toaster />
    </html>
  );
}
