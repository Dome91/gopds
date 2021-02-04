import {Component, Input, isDevMode, OnInit} from '@angular/core';
import {faDownload} from "@fortawesome/free-solid-svg-icons";
import {CatalogEntry} from "../../models/catalog";

@Component({
  selector: 'app-catalog-entry',
  templateUrl: './catalog-entry.component.html',
  styleUrls: ['./catalog-entry.component.sass']
})
export class CatalogEntryComponent implements OnInit {

  @Input() catalogEntry: CatalogEntry;
  downloadURL: string;
  coverURL: string;
  faDownload = faDownload;

  constructor() {
    this.catalogEntry = CatalogEntry.empty();
    this.downloadURL = '';
    this.coverURL = '';
  }

  ngOnInit(): void {
    this.determineDownloadURL();
    this.determineCoverURL();
  }

  private determineDownloadURL() {
    if (isDevMode()) {
      this.downloadURL = `http://localhost:3000/api/v1/catalog/${this.catalogEntry.id}/download`;
    } else {
      this.downloadURL = `api/v1/catalog/${this.catalogEntry.id}/download`;
    }
  }

  private determineCoverURL() {
    if (this.catalogEntry.cover === '') {
      this.coverURL = "assets/card.svg";
    } else if (isDevMode()) {
      this.coverURL = `http://localhost:3000/covers/${this.catalogEntry.cover}`;
    } else {
      this.coverURL = `covers/${this.catalogEntry.cover}`;
    }
  }

}
