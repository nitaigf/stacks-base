export type User = {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'member';
  status: 'active' | 'blocked';
  createdAt: string;
  updatedAt: string;
};

export type AuthPayload = {
  accessToken: string;
  user: User;
};

export type AuthEnvelope = {
  data: AuthPayload;
};

export type UserEnvelope = {
  data: User;
};

export type ApiErrorEnvelope = {
  error: {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
};