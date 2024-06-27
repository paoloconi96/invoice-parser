import { GoogleLogin } from "@react-oauth/google";

const baseApiUrl = "localhost:8000/api/v1";

export default function Login() {
  return (
    <main className="flex min-h-screen flex-col items-center p-24">
      <GoogleLogin
        onSuccess={async credentialResponse => {
          console.log(credentialResponse);

          const request = new Request(`//${baseApiUrl}/login`, {
            body: JSON.stringify(credentialResponse),
            method: "POST",
            credentials: 'include',
          });

          const response = await fetch(request);
        }}
        onError={() => {
          console.log('Login Failed');
        }}
        use_fedcm_for_prompt={true}
      />
    </main>
  );
}
