interface LoginResponse {
  success: boolean;
  message: string;
  data: {
    user: {
      id: string;
      name: string;
      email: string;
      created_at: Date;
      updated_at: Date;
    };
    access_token: string;
    refresh_token: string;
    token_type: string;
    expires_in: number;
  };
  timestamp: string;
}