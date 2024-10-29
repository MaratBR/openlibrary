import LoginForm from "./LoginForm";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { TabsContent } from "@radix-ui/react-tabs";
import SignUpForm from "./SignUpForm";
import useAuthState from "../state";
import { Navigate } from "react-router";

export default function LoginPage() {
  const state = useAuthState();

  if (state.user) {
    return <Navigate to="/home" />;
  }

  return (
    <div className="grid h-screen w-screen grid-cols-2">
      <div>123</div>
      <div className="p-10 flex items-center">
        <div className="p-6">
          <Tabs defaultValue="signin">
            <TabsList>
              <TabsTrigger value="signin">Sign in</TabsTrigger>
              <TabsTrigger value="signup">Sign up</TabsTrigger>
            </TabsList>
            <TabsContent value="signin">
              <div className="min-h-[300px]">
                <LoginForm />
              </div>
            </TabsContent>
            <TabsContent value="signup">
              <div className="min-h-[300px]">
                <SignUpForm />
              </div>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  );
}
