import Link from "next/link";

import { Post } from "@/app/_components/post";
import { auth } from "@/server/auth";
import { api, HydrateClient } from "@/trpc/server";

export default async function Home() {
  const session = await auth();

  if (session?.user) {
    void api.post.getLatest.prefetch();
  }

  return (
    <HydrateClient>
      <main className="flex min-h-screen flex-col items-center bg-gradient-to-b from-[#2e026d] to-[#15162c] text-white">
        <div className="w-full border-b border-gray-700">
          <div className="container mx-auto flex items-center justify-end p-4">
            <div className="flex items-center space-x-4">
              {session && (
                <span className="text-lg font-medium">
                  {session.user?.name}
                </span>
              )}
              <Link
                href={session ? "/api/auth/signout" : "/api/auth/signin"}
                className="rounded-full bg-white/10 px-4 py-2 text-sm font-semibold no-underline transition hover:bg-white/20"
              >
                {session ? "Выйти" : "Войти"}
              </Link>
            </div>
          </div>
        </div>

        <div className="container flex flex-grow flex-col items-center justify-center px-4 py-8">
          {session?.user && <Post />}
        </div>
      </main>
    </HydrateClient>
  );
}
