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
  downloadLink: string;

  faDownload = faDownload;

  constructor() {
    this.catalogEntry = new CatalogEntry('', '', true);
    this.downloadLink = '';
  }

  ngOnInit(): void {
    if (isDevMode()) {
      this.downloadLink = `http://localhost:3000/api/v1/catalog/${this.catalogEntry.id}/download`;
    } else {
      this.downloadLink = `api/v1/catalog/${this.catalogEntry.id}/download`;
    }
  }

}
