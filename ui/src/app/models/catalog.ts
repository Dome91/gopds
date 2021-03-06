export class CatalogEntriesInPage {
  constructor(public total: number, public catalogEntries: CatalogEntry[]) {
  }
}

export class CatalogEntry {
  constructor(public id: string, public name: string, public isDirectory: boolean, public type: string, public cover: string) {
  }

  public static empty(): CatalogEntry {
    return new CatalogEntry('', '', true, "", "");
  }
}
