import { httpClient } from "../common/api";

type SignInRequest = {
  username: string;
  password: string;
};

export async function httpSignIn(req: SignInRequest): Promise<void> {
  await httpClient.post("/api/auth/signin", { json: req });
}

type SignUpRequest = {
  username: string;
  password: string;
};

export async function httpSignUp(req: SignUpRequest): Promise<void> {
  await httpClient.post("/api/auth/signup", { json: req });
}
