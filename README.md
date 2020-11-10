# mostlyzero: a tool to open a block device and search for content

I do some odd things with hard drives. Thinks like:
* passing large swaths of them directly to VM's without using a storage pool
* skipping partition tables wholesale for certain obstinate kernels/firmware that behave oddly.
* luks directly of the block device (no partition table)
* liberal use of dd and losetup of devices with indeterminate offsets into the device
* occasional power-loss mid operation

As a byproduct, I end up with disks with unexpected bytes in unexpected places. Initially, this program will:
1) open a block device directly
2) do a sort of binary search of the drive for non-zero data
3) exit if nonzero data is found, and report the location

I can think of a few other addons which would be nice:
* Sometimes I have run things like dban on old drives. Detecting dban patterns would be useful, so we can stop looking at that part of the disk
* spread searching. If the program identified data, searching for the range of blocks on the disk where it found the data could be useful.
* Entropy detection. If data is found, doing some entropy detection of the data to find out if it is likely  encrypted/`dd if=/dev/urandon`/etc.
