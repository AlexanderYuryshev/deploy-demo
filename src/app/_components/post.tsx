"use client";

import { useState } from "react";

import { api } from "@/trpc/react";

export function Post() {
  const [latestPost] = api.post.getLatest.useSuspenseQuery();

  const utils = api.useUtils();
  const [name, setName] = useState("");
  const [content, setContent] = useState("");
   const [generationStyle, setGenerationStyle] = useState("informal");

  const createPost = api.post.create.useMutation({
    onSuccess: async () => {
      await utils.post.invalidate();
      setName("");
    },
  });
  const generatePost = api.generate.generatePost.useMutation({
    onSuccess: async (response) => {
      setContent(response);
    },
  });

  return (
    <div className="w-full">
      {latestPost ? (
        <p className="truncate">Ваш последний пост: {latestPost.name}</p>
      ) : (
        <p>У вас еще нет постов</p>
      )}
      <form
        onSubmit={(e) => {
          e.preventDefault();
          createPost.mutate({ name, content });
        }}
        className="flex flex-col gap-2"
      >
        <div className="flex flex-grow flex-col p-4">
          <div className="mb-6 w-full">
            <input
              type="text"
              placeholder="Название поста"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full rounded-lg bg-white/10 px-6 py-4 text-center text-lg text-white"
            />
          </div>

          <div className="mb-6 w-full flex-grow">
            <textarea
              placeholder="Текст поста"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              className="h-full w-full resize-y rounded-lg bg-white/10 px-6 py-4 text-white"
            />
          </div>

          <div className="flex w-full space-x-4 justify-center">
            <select
              value={generationStyle}
              onChange={(e) => setGenerationStyle(e.target.value)}
              className="cursor-pointer appearance-none rounded-full bg-gray-500 px-4 py-3 text-center font-semibold"
            >
              <option value="informal">Неформальный</option>
              <option value="professional">Профессиональный</option>
              <option value="humorous">Юмористический</option>
            </select>
            <button
              type="button"
              className="rounded-full bg-purple-500 px-4 py-3 font-semibold hover:bg-purple-600 disabled:opacity-50"
              disabled={createPost.isPending || !name}
              onClick={() => name && generatePost.mutate({ name, style: generationStyle as "informal" | "professional" | "humorous" })}
            >
              Сгенерировать
            </button>

            <button
              type="submit"
              className="rounded-full bg-green-500 px-4 py-3 font-semibold hover:bg-green-600 disabled:opacity-50"
              disabled={createPost.isPending}
            >
              {createPost.isPending ? "Загрузка..." : "Опубликовать"}
            </button>
          </div>
        </div>
      </form>
    </div>
  );
}
