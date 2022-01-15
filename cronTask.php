<?php
$msgs = [
	[
		'date' => date('Y-m-d'),
		'source_id' => 'aaaaaaaaaaa',
		'source_type' => 'api',// api
		'urls' => [
			[
				'url' => 'https://weibo.com/ajax/statuses/mymblog',
				'params' => [
					'page' => '1',
					'feature' => '0',
					'uid' => '1270577247'
				],
				'fields' => [],
				'|' => [
					'url' => '',
					'params' => [
						'foo' => 'bar'	
					],
					'fields' => [],
					'|' => [
						'url' => '',
						'params' => [
							'foo' => 'bar'
						],
						'fields' => [],
						'...'

					]
				]
			],
			[
				'...'
			],
		]
	],
	[
		'date' => date('Y-m-d'),
		'source_id' => 'bbbbbbbbbbb',
		'source_type' => 'file',// multimedia file
		'paths' => [
			[
				'account' => 'root',
				'pwd' => 'xxxxxxxxx', //decrypt
				'host' => 'tcp://xxxxxxxxxxxxx',
				'uri' => '/root',
				'suffix' => ['mp4', 'doc'],
				// 'saveTo' => '/tmp/collect/csv/s323sasf23a46456qe2sdfd32x/20220108.doc',
				'saveTo' => '/tmp/collect'
			],
		]
	],
	[
		'date' => date('Y-m-d'),
		'source_id' => 'eeeeeeeee',
		'source_type' => 'csv',// multimedia filer
		'paths' => [
			[
				'account' => 'root',
				'pwd' => 'xxxxxxxxxxxxx', //decrypt
				'host' => 'tcp://xxxxxxxxxxxxxx',
				'uri' => "/root/install.sh",
				// 'saveTo' => '/tmp/collect/csv/s323sasf23a46456qe2sdfd32x/20220108.doc',
				'saveTo' => '/tmp/collect'
			],
		]
	],
	[
		'date' => date('Y-m-d'),
		'source_id' => 'cccccccccc',
		'source_type' => 'mongodb',// direct mongodb
		'queries' => [
			[
				'query' => 'xxxxxxxxxxxxxxxxxxxx',
				'params' => [
					'page' => '1',
					'feature' => '0',
					'uid' => '1270577247'
				],
			],
		]
	],
	[
		'date' => date('Y-m-d'),
		'source_id' => 'ddddddddddd',
		'source_type' => 'mysql',// direct database
		'queries' => [
			[
				'query' => 'https://weibo.com/ajax/statuses/mymblog',
				'params' => [
					'page' => '1',
					'feature' => '0',
					'uid' => '1270577247'
				],
			],
		]
	],
		
];

$msgs = [];
$msgs[] = [
	'date' => date('Y-m-d'),
	'source_id' => 's323sasf23a46456qe2sdfd32x',
	'source_type' => 'csv',// multimedia file
	'paths' => [
		// [
		// 	'account' => 'root',
		// 	'pwd' => 'xxxxxx', //decrypt
		// 	'host' => 'tcp://xxxxxxxxxxx',
		// 	'uri' => "/root/testfile",
		// 	'suffix' => [],
		// 	'saveTo' => '/tmp/collect/csv/s323sasf23a46456qe2sdfd32x/20220108.csv',
		// ],
		[
			'account' => 'root',
			'pwd' => 'xxxxxxxx', //decrypt
			'host' => 'tcp://xxxxxxxxxx:22',
			'uri' => '/root/test',
			'suffix' => ['mp4', 'doc'],
			// 'saveTo' => '/tmp/collect',
			'saveTo' => '/tmp/collect/multimedia'
		],
	]
];

$socket = socket_create(AF_UNIX, SOCK_STREAM, 0);
socket_connect($socket, '/tmp/collector.sock');

foreach($msgs as $msg) {
	$msg = json_encode($msg);
    socket_send($socket, $msg, strlen($msg), 0);
    $response = socket_read($socket, 1024);
	sleep(1);
    if($response != 'ok') {
		//todo log error
	}
}

socket_close($socket);






