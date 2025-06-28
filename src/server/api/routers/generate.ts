import { z } from "zod";
import { ChatOpenAI } from "@langchain/openai";

import { env } from "@/env";
import { createTRPCRouter, protectedProcedure } from "@/server/api/trpc";

const baseURL = "https://models.github.ai/inference";
const model = "openai/gpt-4.1";

const systemPrompt = `
Ты — блогер. Напиши пост для социальной сети на заданную тему.

Требования:
1.  Начни с привлекающего внимание введения (избегай шаблонов вроде "Привет, подписчики!").
2.  Раскрой тему через неочевидный ракурс, свежую метафору, личный опыт с изюминкой или малоизвестный факт. Избегай поверхностных советов, общеизвестных истин, клише и шаблонных фраз.
3.  Вызови искренний эмоциональный отклик (любопытство, удивление, вдохновение и т.д.). Будь аутентичным.
4.  Заверши пост сильным призывом к действию, провокационным вопросом или запоминающейся фразой (не просто "А что вы думаете?").
5.  Объем: 3-5 абзацев.
6.  Запрещено: хэштеги, упоминания (@), пересказ общедоступной информации.

Пиши на русском языке.
`.trim();

const createSafeChatModel = () => {
  if (!env.OPENAI_API_KEY) {
    throw new Error("OPENAI_API_KEY is not configured");
  }

  return new ChatOpenAI({
    model,
    apiKey: env.OPENAI_API_KEY,
    configuration: { baseURL },
    timeout: 30_000,
    maxRetries: 1,
  });
};

export const generateRouter = createTRPCRouter({
  generatePost: protectedProcedure
    .input(
      z.object({
        name: z
          .string()
          .min(5, "Тема поста должна содержать минимум 5 символов"),
        style: z
          .enum(["informal", "professional", "humorous"])
          .optional()
          .default("informal"),
      }),
    )
    .mutation(async ({ input }) => {
      const { name, style } = input;
      try {
        const openAImodel = createSafeChatModel();

        const styleDescription = {
          informal: "неформальном разговорном стиле",
          professional: "профессиональном деловом стиле",
          humorous: "юмористическом и легком стиле",
        }[style];

        const response = await openAImodel.invoke([
          { role: "system", content: systemPrompt },
          {
            role: "user",
            content: `Напиши пост на тему: "${name}" в ${styleDescription}.`,
          },
        ]);

        return response.text;
      } catch (error) {
        console.error("Ошибка генерации поста:", error);

        let errorMessage = "Не удалось сгенерировать пост";
        if (error instanceof Error) {
          errorMessage += `: ${error.message}`;
        }

        throw new Error(errorMessage);
      }
    }),
});
