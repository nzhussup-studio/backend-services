import type { Meta, StoryObj } from "../kc.gen";
import { createKcPageStory } from "./KcPageStory";

const { KcPageStory } = createKcPageStory({ pageId: "login.ftl" });

const meta = {
  title: "Login/Sign in",
  component: KcPageStory
} satisfies Meta<typeof KcPageStory>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithSocialProviders: Story = {
  args: {
    kcContext: {
      social: {
        displayInfo: true,
        providers: [
          { loginUrl: "#", alias: "google", providerId: "google", displayName: "Google" },
          { loginUrl: "#", alias: "github", providerId: "github", displayName: "GitHub" }
        ]
      }
    }
  }
};

export const InvalidCredentials: Story = {
  args: {
    kcContext: {
      message: { type: "error", summary: "Invalid username or password." }
    }
  }
};
