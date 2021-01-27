export class User {
  constructor(public username: string,
              public password: string,
              public role: Role) {
  }
}

export enum Role {
  ADMIN = 'ADMIN',
  USER = 'USER'
}
