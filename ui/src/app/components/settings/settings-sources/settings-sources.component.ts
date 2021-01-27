import {Component, OnInit} from '@angular/core';
import {Source} from "../../../models/source";
import {SourceService} from "../../../services/source.service";
import {ToastrService} from "ngx-toastr";
import {faTrash, faSync,} from "@fortawesome/free-solid-svg-icons";
import {NgbModal} from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: 'app-settings-sources',
  templateUrl: './settings-sources.component.html',
  styleUrls: ['./settings-sources.component.sass']
})
export class SettingsSourcesComponent implements OnInit {

  faTrash = faTrash;
  faSync = faSync;

  sources: Source[];
  isSyncing: Map<string, boolean>;

  constructor(private sourceService: SourceService, private toastrService: ToastrService, private modalService: NgbModal) {
    this.sources = [];
    this.isSyncing = new Map<string, boolean>();
  }

  ngOnInit(): void {
    this.fetchSources();
  }

  sync(id: string) {
    this.isSyncing.set(id, true);
    this.sourceService.synchronize(id).subscribe(
      {
        next: () => this.toastrService.success('Started synchronization. This may take a while.'),
        complete: () => this.isSyncing.set(id, false)
      });
  }

  delete(id: string) {
    this.sourceService.delete(id).subscribe(
      () => this.fetchSources(),
    );
  }

  private fetchSources() {
    this.sourceService.fetchAll()
      .subscribe(sources => {
        this.sources = sources
        this.isSyncing.clear();
        sources.forEach(source => this.isSyncing.set(source.id, false))
      });
  }

  shouldSpin(id: string): boolean {
    const isSyncing = this.isSyncing.get(id);
    if (isSyncing === undefined) {
      return false;
    }

    return isSyncing;
  }
}
