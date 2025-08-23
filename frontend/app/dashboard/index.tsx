import { useAuthStore } from "@/src/stores/authStore";
import {  Redirect } from "expo-router"; 

export default function Dashboard() {
  const { user } = useAuthStore()
  console.log('user dashboard', user);
  if (!user) {
    return <Redirect href="/auth/login" />
  }
  if (user) {
    return <Redirect href="/dashboard/tabs/overview" />
  }
  return <Redirect href="/auth/login" />
}