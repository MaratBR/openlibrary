import { Card } from "@/components/ui/card";
import LoginForm from "./LoginForm";

export default function LoginPage() {
  return (
    <div className="grid h-screen w-screen grid-cols-2">
      <div>123</div>
      <div className="p-10 flex items-center">
        <div className="p-6">
          <LoginForm />
        </div>
      </div>
    </div>
  );
}
