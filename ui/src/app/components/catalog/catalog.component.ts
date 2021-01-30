import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Params, Router} from "@angular/router";
import {CatalogEntry} from "../../models/catalog";
import {CatalogService} from "../../services/catalog.service";

@Component({
  selector: 'app-catalog',
  templateUrl: './catalog.component.html',
  styleUrls: ['./catalog.component.sass']
})
export class CatalogComponent implements OnInit {

  catalogEntries: CatalogEntry[];
  id: string;
  total: number;
  page: number;
  pageSize: number;

  constructor(private route: ActivatedRoute,
              private router: Router,
              private catalogService: CatalogService) {
    this.catalogEntries = [];
    this.total = 0;
    this.page = 0;
    this.pageSize = 24;
    this.id = '';
  }

  ngOnInit(): void {
    this.route.queryParams.subscribe((params: Params) => {
      this.id = params['id'];
      if (this.id === undefined) {
        this.id = '';
      }
      this.catalogEntries = [];
      this.total = 0;
      this.page = 0;
      this.pageSize = 24;
      this.fetchCatalogEntriesInPage();
    })
  }

  onScroll() {
    const isNotLastPage = this.total > (this.page + 1) * this.pageSize
    if (isNotLastPage) {
      this.page += 1;
      this.fetchCatalogEntriesInPage();
    }
  }

  private fetchCatalogEntriesInPage() {
    this.catalogService.fetchInPage(this.page, this.pageSize, this.id).subscribe(
      response => {
        this.catalogEntries.push(...response.catalogEntries);
        this.total = response.total;
      }
    );
  }

  navigateToCatalogEntry(id: string) {
    this.router.navigateByUrl(`catalog?id=${id}`, {
    });
  }
}
